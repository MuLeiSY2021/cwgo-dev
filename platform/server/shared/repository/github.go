/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package repository

import (
	"context"
	"fmt"
	"github.com/google/go-github/v53/github"
	"io/ioutil"
	"net/http"
)

type GitHubApi struct {
	Client map[int64]*github.Client
}

func (g *GitHubApi) GetFile(repoId int64, owner, repoName, filePath, ref string) (*File, error) {
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}
	fileContent, _, err := g.Client[repoId].Repositories.DownloadContents(context.Background(), owner, repoName, filePath, opts)
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	content, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return nil, err
	}

	return &File{
		Name:    filePath,
		Content: content,
	}, nil
}

func (g *GitHubApi) PushFilesToRepository(files map[string][]byte, repoId int64, owner, repoName, branch, commitMessage string) error {
	// Get a reference to the default branch
	ref, _, err := g.Client[repoId].Git.GetRef(context.Background(), owner, repoName, "refs/heads/"+branch)
	if err != nil {
		return err
	}

	// Create a new Tree object for the file to be pushed
	var treeEntries []*github.TreeEntry
	for filePath, content := range files {
		treeEntries = append(treeEntries, &github.TreeEntry{
			Path:    github.String(filePath),
			Content: github.String(string(content)),
			Mode:    github.String("100644"),
		})
	}
	newTree, _, err := g.Client[repoId].Git.CreateTree(context.Background(), owner, repoName, *ref.Object.SHA, treeEntries)
	if err != nil {
		return err
	}

	// Create a new commit object, using the new tree as its foundation
	newCommit, _, err := g.Client[repoId].Git.CreateCommit(context.Background(), owner, repoName, &github.Commit{
		Message: github.String(commitMessage),
		Tree:    newTree,
	})
	if err != nil {
		return err
	}

	// Update branch references to point to new submissions
	_, _, err = g.Client[repoId].Git.UpdateRef(context.Background(), owner, repoName, &github.Reference{
		Ref: github.String("refs/heads/" + branch),
		Object: &github.GitObject{
			SHA:  newCommit.SHA,
			Type: github.String("commit"),
		},
	}, true)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitHubApi) GetRepositoryArchive(repoId int64, owner, repoName, format, ref string) ([]byte, error) {
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	archiveLink, _, err := g.Client[repoId].Repositories.GetArchiveLink(context.Background(), owner, repoName, github.ArchiveFormat(format), opts, false)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(archiveLink.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch archive: %s", resp.Status)
	}

	archiveData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return archiveData, nil
}

func (g *GitHubApi) GetLatestCommitHash(repoId int64, owner, repoName, filePath, ref string) (string, error) {
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	fileContent, _, _, err := g.Client[repoId].Repositories.GetContents(context.Background(), owner, repoName, filePath, opts)
	if err != nil {
		return "", err
	}

	return *fileContent.SHA, nil
}
